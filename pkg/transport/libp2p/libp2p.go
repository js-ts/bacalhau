package libp2p

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/filecoin-project/bacalhau/pkg/system"
	"github.com/filecoin-project/bacalhau/pkg/transport"
	"github.com/filecoin-project/bacalhau/pkg/types"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog/log"
)

const JOB_EVENT_CHANNEL = "bacalhau-job-event"

type Transport struct {
	// Cleanup manager for resource teardown on exit:
	cm *system.CleanupManager

	// Writer we emit events through.
	genericTransport     *transport.GenericTransport
	Host                 host.Host
	Port                 int
	PubSub               *pubsub.PubSub
	JobEventTopic        *pubsub.Topic
	JobEventSubscription *pubsub.Subscription
}

func getConfigPath() string {
	suffix := "/.bacalhau"
	env := os.Getenv("BACALHAU_PATH")
	var d string
	if env == "" {
		// e.g. /home/francesca/.bacalhau
		dirname, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		d = dirname + suffix
	} else {
		// e.g. /data/.bacalhau
		d = env + suffix
	}
	// create dir if not exists
	if err := os.MkdirAll(d, 0700); err != nil {
		panic(err)
	}
	return d
}

func makeLibp2pHost(port int) (host.Host, error) {
	configPath := getConfigPath()

	// We include the port in the filename so that in devstack multiple nodes
	// running on the same host get different identities
	privKeyPath := fmt.Sprintf("%s/private_key.%d", configPath, port)

	if _, err := os.Stat(privKeyPath); errors.Is(err, os.ErrNotExist) {
		// Private key does not exist - create and write it

		// Creates a new RSA key pair for this host.
		prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
		if err != nil {
			log.Error().Err(err)
			return nil, err
		}

		keyOut, err := os.OpenFile(privKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to open key.pem for writing: %v", err)
		}
		privBytes, err := crypto.MarshalPrivateKey(prvKey)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal private key: %v", err)
		}
		// base64 encode privBytes
		b64 := base64.StdEncoding.EncodeToString(privBytes)
		_, err = keyOut.Write([]byte(b64 + "\n"))
		if err != nil {
			return nil, fmt.Errorf("failed to write to key file: %v", err)
		}
		if err := keyOut.Close(); err != nil {
			return nil, fmt.Errorf("error closing key file: %v", err)
		}
		log.Printf("wrote %s", privKeyPath)
	}

	// Now that we've ensured the private key is written to disk, read it! This
	// ensures that loading it works even in the case where we've just created
	// it.

	// read the private key
	keyBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}
	// base64 decode keyBytes
	b64, err := base64.StdEncoding.DecodeString(string(keyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}
	// parse the private key
	prvKey, err := crypto.UnmarshalPrivateKey(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	return libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
}

func NewTransport(cm *system.CleanupManager, port int) (
	*Transport, error) {

	host, err := makeLibp2pHost(port)
	if err != nil {
		return nil, err
	}

	// libp2p uses contexts for lifecycle management
	ctx, cancel := context.WithCancel(context.Background())
	cm.RegisterCallback(func() error {
		cancel()
		return ctx.Err()
	})

	pubsub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}

	jobEventTopic, err := pubsub.Join(JOB_EVENT_CHANNEL)
	if err != nil {
		return nil, err
	}

	jobEventSubscription, err := jobEventTopic.Subscribe()
	if err != nil {
		return nil, err
	}

	libp2pTransport := &Transport{
		cm:                   cm,
		Host:                 host,
		Port:                 port,
		PubSub:               pubsub,
		JobEventTopic:        jobEventTopic,
		JobEventSubscription: jobEventSubscription,
	}

	// setup the event writer
	libp2pTransport.genericTransport = transport.NewGenericTransport(
		host.ID().String(),
		func(ctx context.Context, event *types.JobEvent) error {
			return libp2pTransport.writeJobEvent(ctx, event)
		},
	)

	return libp2pTransport, nil
}

/////////////////////////////////////////////////////////////
/// LIFECYCLE
/////////////////////////////////////////////////////////////

func (t *Transport) HostID(ctx context.Context) (string, error) {
	return t.Host.ID().String(), nil
}

func (t *Transport) Start(ctx context.Context) error {
	if len(t.genericTransport.SubscribeFuncs) <= 0 {
		panic("Programming error: no subscribe func, please call Subscribe immediately after constructing interface")
	}

	go t.readLoopJobEvents(ctx)
	log.Debug().Msg("Libp2p transport has started")

	t.cm.RegisterCallback(func() error {
		t.Host.Close()
		log.Debug().Msg("Libp2p transport has stopped")
		return nil
	})

	log.Debug().Msg("libp2p transport is starting...")
	t.readLoopJobEvents(ctx) // blocking

	return nil
}

/////////////////////////////////////////////////////////////
/// READ OPERATIONS
/////////////////////////////////////////////////////////////

func (t *Transport) List(ctx context.Context) (
	types.ListResponse, error) {

	return t.genericTransport.List(ctx)
}

func (t *Transport) Get(ctx context.Context, id string) (*types.Job, error) {
	return t.genericTransport.Get(ctx, id)
}

func (t *Transport) Subscribe(ctx context.Context, fn func(
	jobEvent *types.JobEvent, job *types.Job)) {

	t.genericTransport.Subscribe(ctx, fn)
}

/////////////////////////////////////////////////////////////
/// WRITE OPERATIONS - "CLIENT" / REQUESTER
/////////////////////////////////////////////////////////////

func (t *Transport) SubmitJob(ctx context.Context, spec *types.JobSpec,
	deal *types.JobDeal) (*types.Job, error) {

	return t.genericTransport.SubmitJob(ctx, spec, deal)
}

func (t *Transport) UpdateDeal(ctx context.Context, jobID string,
	deal *types.JobDeal) error {

	return t.genericTransport.UpdateDeal(ctx, jobID, deal)
}

func (t *Transport) CancelJob(ctx context.Context, jobID string) error {
	return nil
}

func (t *Transport) AcceptJobBid(ctx context.Context, jobID,
	nodeID string) error {

	return t.genericTransport.AcceptJobBid(ctx, jobID, nodeID)
}

func (t *Transport) RejectJobBid(ctx context.Context, jobID, nodeID,
	message string) error {

	return t.genericTransport.RejectJobBid(ctx, jobID, nodeID, message)
}

/////////////////////////////////////////////////////////////
/// WRITE OPERATIONS - "SERVER" / COMPUTE NODE
/////////////////////////////////////////////////////////////

func (t *Transport) BidJob(ctx context.Context, jobID string) error {
	return t.genericTransport.BidJob(ctx, jobID)
}

func (t *Transport) SubmitResult(ctx context.Context, jobID, status,
	resultsID string) error {

	return t.genericTransport.SubmitResult(ctx, jobID, status, resultsID)
}

func (t *Transport) ErrorJob(ctx context.Context, jobID, status string) error {
	return t.genericTransport.ErrorJob(ctx, jobID, status)
}

// this is when the requester node needs to error the status for a node
// for example - results have been given by the compute node
// and in checking the results, the requester node came across some kind of error
// we need to flag that error against the node that submitted the results
// (but we are the requester node) - so we need this util function
func (t *Transport) ErrorJobForNode(ctx context.Context, jobID, nodeID,
	status string) error {

	return t.genericTransport.ErrorJobForNode(ctx, jobID, nodeID, status)
}

/////////////////////////////////////////////////////////////
/// INTERNAL IMPLEMENTATION
/////////////////////////////////////////////////////////////

func (t *Transport) Connect(ctx context.Context, peerConnect string) error {
	if peerConnect == "" {
		return nil
	}
	maddr, err := multiaddr.NewMultiaddr(peerConnect)
	if err != nil {
		return err
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}

	t.Host.Peerstore().AddAddrs(
		info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	return t.Host.Connect(ctx, *info)
}

func (t *Transport) writeJobEvent(ctx context.Context, event *types.JobEvent) error {
	bs, err := json.Marshal(event)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Sending event: %s", string(bs))
	return t.JobEventTopic.Publish(ctx, bs)
}

func (t *Transport) readLoopJobEvents(ctx context.Context) {
	for {
		msg, err := t.JobEventSubscription.Next(ctx)
		if err != nil {
			return
		}

		jobEvent := new(types.JobEvent)
		err = json.Unmarshal(msg.Data, jobEvent)
		if err != nil {
			continue
		}

		t.genericTransport.BroadcastEvent(jobEvent)
	}
}

// Static check to ensure that Transport implements Transport:
var _ transport.Transport = (*Transport)(nil)

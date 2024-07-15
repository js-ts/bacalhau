# coding: utf-8

"""
    Bacalhau API

    This page is the reference of the Bacalhau REST API. Project docs are available at https://docs.bacalhau.org/. Find more information about Bacalhau at https://github.com/bacalhau-project/bacalhau.  # noqa: E501

    OpenAPI spec version: ${VERSION}
    Contact: team@bacalhau.org
    Generated by: https://github.com/swagger-api/swagger-codegen.git
"""

import pprint
import re  # noqa: F401

import six

class InputSource(object):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """
    """
    Attributes:
      swagger_types (dict): The key is attribute name
                            and the value is attribute type.
      attribute_map (dict): The key is attribute name
                            and the value is json key in definition.
    """
    swagger_types = {
        'alias': 'str',
        'source': 'AllOfInputSourceSource',
        'target': 'str'
    }

    attribute_map = {
        'alias': 'Alias',
        'source': 'Source',
        'target': 'Target'
    }

    def __init__(self, alias=None, source=None, target=None):  # noqa: E501
        """InputSource - a model defined in Swagger"""  # noqa: E501
        self._alias = None
        self._source = None
        self._target = None
        self.discriminator = None
        if alias is not None:
            self.alias = alias
        if source is not None:
            self.source = source
        if target is not None:
            self.target = target

    @property
    def alias(self):
        """Gets the alias of this InputSource.  # noqa: E501

        Alias is an optional reference to this input source that can be used for dynamic linking to this input. (e.g. dynamic import in wasm by alias)  # noqa: E501

        :return: The alias of this InputSource.  # noqa: E501
        :rtype: str
        """
        return self._alias

    @alias.setter
    def alias(self, alias):
        """Sets the alias of this InputSource.

        Alias is an optional reference to this input source that can be used for dynamic linking to this input. (e.g. dynamic import in wasm by alias)  # noqa: E501

        :param alias: The alias of this InputSource.  # noqa: E501
        :type: str
        """

        self._alias = alias

    @property
    def source(self):
        """Gets the source of this InputSource.  # noqa: E501

        Source is the source of the artifact to be downloaded, e.g a URL, S3 bucket, etc.  # noqa: E501

        :return: The source of this InputSource.  # noqa: E501
        :rtype: AllOfInputSourceSource
        """
        return self._source

    @source.setter
    def source(self, source):
        """Sets the source of this InputSource.

        Source is the source of the artifact to be downloaded, e.g a URL, S3 bucket, etc.  # noqa: E501

        :param source: The source of this InputSource.  # noqa: E501
        :type: AllOfInputSourceSource
        """

        self._source = source

    @property
    def target(self):
        """Gets the target of this InputSource.  # noqa: E501

        Target is the path where the artifact should be mounted on  # noqa: E501

        :return: The target of this InputSource.  # noqa: E501
        :rtype: str
        """
        return self._target

    @target.setter
    def target(self, target):
        """Sets the target of this InputSource.

        Target is the path where the artifact should be mounted on  # noqa: E501

        :param target: The target of this InputSource.  # noqa: E501
        :type: str
        """

        self._target = target

    def to_dict(self):
        """Returns the model properties as a dict"""
        result = {}

        for attr, _ in six.iteritems(self.swagger_types):
            value = getattr(self, attr)
            if isinstance(value, list):
                result[attr] = list(map(
                    lambda x: x.to_dict() if hasattr(x, "to_dict") else x,
                    value
                ))
            elif hasattr(value, "to_dict"):
                result[attr] = value.to_dict()
            elif isinstance(value, dict):
                result[attr] = dict(map(
                    lambda item: (item[0], item[1].to_dict())
                    if hasattr(item[1], "to_dict") else item,
                    value.items()
                ))
            else:
                result[attr] = value
        if issubclass(InputSource, dict):
            for key, value in self.items():
                result[key] = value

        return result

    def to_str(self):
        """Returns the string representation of the model"""
        return pprint.pformat(self.to_dict())

    def __repr__(self):
        """For `print` and `pprint`"""
        return self.to_str()

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        if not isinstance(other, InputSource):
            return False

        return self.__dict__ == other.__dict__

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        return not self == other
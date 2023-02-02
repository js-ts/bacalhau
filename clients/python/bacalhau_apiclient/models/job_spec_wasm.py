# coding: utf-8

"""
    Bacalhau API

    This page is the reference of the Bacalhau REST API. Project docs are available at https://docs.bacalhau.org/. Find more information about Bacalhau at https://github.com/filecoin-project/bacalhau.  # noqa: E501

    OpenAPI spec version: 0.3.18.post4
    Contact: team@bacalhau.org
    Generated by: https://github.com/swagger-api/swagger-codegen.git
"""


import pprint
import re  # noqa: F401

import six

from bacalhau_apiclient.configuration import Configuration


class JobSpecWasm(object):
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
        'entry_module': 'JobSpecWasmEntryModule',
        'entry_point': 'str',
        'environment_variables': 'dict(str, str)',
        'import_modules': 'list[StorageSpec]',
        'parameters': 'list[str]'
    }

    attribute_map = {
        'entry_module': 'EntryModule',
        'entry_point': 'EntryPoint',
        'environment_variables': 'EnvironmentVariables',
        'import_modules': 'ImportModules',
        'parameters': 'Parameters'
    }

    def __init__(self, entry_module=None, entry_point=None, environment_variables=None, import_modules=None, parameters=None, _configuration=None):  # noqa: E501
        """JobSpecWasm - a model defined in Swagger"""  # noqa: E501
        if _configuration is None:
            _configuration = Configuration()
        self._configuration = _configuration

        self._entry_module = None
        self._entry_point = None
        self._environment_variables = None
        self._import_modules = None
        self._parameters = None
        self.discriminator = None

        if entry_module is not None:
            self.entry_module = entry_module
        if entry_point is not None:
            self.entry_point = entry_point
        if environment_variables is not None:
            self.environment_variables = environment_variables
        if import_modules is not None:
            self.import_modules = import_modules
        if parameters is not None:
            self.parameters = parameters

    @property
    def entry_module(self):
        """Gets the entry_module of this JobSpecWasm.  # noqa: E501


        :return: The entry_module of this JobSpecWasm.  # noqa: E501
        :rtype: JobSpecWasmEntryModule
        """
        return self._entry_module

    @entry_module.setter
    def entry_module(self, entry_module):
        """Sets the entry_module of this JobSpecWasm.


        :param entry_module: The entry_module of this JobSpecWasm.  # noqa: E501
        :type: JobSpecWasmEntryModule
        """

        self._entry_module = entry_module

    @property
    def entry_point(self):
        """Gets the entry_point of this JobSpecWasm.  # noqa: E501

        The name of the function in the EntryModule to call to run the job. For WASI jobs, this will always be `_start`, but jobs can choose to call other WASM functions instead. The EntryPoint must be a zero-parameter zero-result function.  # noqa: E501

        :return: The entry_point of this JobSpecWasm.  # noqa: E501
        :rtype: str
        """
        return self._entry_point

    @entry_point.setter
    def entry_point(self, entry_point):
        """Sets the entry_point of this JobSpecWasm.

        The name of the function in the EntryModule to call to run the job. For WASI jobs, this will always be `_start`, but jobs can choose to call other WASM functions instead. The EntryPoint must be a zero-parameter zero-result function.  # noqa: E501

        :param entry_point: The entry_point of this JobSpecWasm.  # noqa: E501
        :type: str
        """

        self._entry_point = entry_point

    @property
    def environment_variables(self):
        """Gets the environment_variables of this JobSpecWasm.  # noqa: E501

        The variables available in the environment of the running program.  # noqa: E501

        :return: The environment_variables of this JobSpecWasm.  # noqa: E501
        :rtype: dict(str, str)
        """
        return self._environment_variables

    @environment_variables.setter
    def environment_variables(self, environment_variables):
        """Sets the environment_variables of this JobSpecWasm.

        The variables available in the environment of the running program.  # noqa: E501

        :param environment_variables: The environment_variables of this JobSpecWasm.  # noqa: E501
        :type: dict(str, str)
        """

        self._environment_variables = environment_variables

    @property
    def import_modules(self):
        """Gets the import_modules of this JobSpecWasm.  # noqa: E501

        TODO #880: Other WASM modules whose exports will be available as imports to the EntryModule.  # noqa: E501

        :return: The import_modules of this JobSpecWasm.  # noqa: E501
        :rtype: list[StorageSpec]
        """
        return self._import_modules

    @import_modules.setter
    def import_modules(self, import_modules):
        """Sets the import_modules of this JobSpecWasm.

        TODO #880: Other WASM modules whose exports will be available as imports to the EntryModule.  # noqa: E501

        :param import_modules: The import_modules of this JobSpecWasm.  # noqa: E501
        :type: list[StorageSpec]
        """

        self._import_modules = import_modules

    @property
    def parameters(self):
        """Gets the parameters of this JobSpecWasm.  # noqa: E501

        The arguments supplied to the program (i.e. as ARGV).  # noqa: E501

        :return: The parameters of this JobSpecWasm.  # noqa: E501
        :rtype: list[str]
        """
        return self._parameters

    @parameters.setter
    def parameters(self, parameters):
        """Sets the parameters of this JobSpecWasm.

        The arguments supplied to the program (i.e. as ARGV).  # noqa: E501

        :param parameters: The parameters of this JobSpecWasm.  # noqa: E501
        :type: list[str]
        """

        self._parameters = parameters

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
        if issubclass(JobSpecWasm, dict):
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
        if not isinstance(other, JobSpecWasm):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, JobSpecWasm):
            return True

        return self.to_dict() != other.to_dict()
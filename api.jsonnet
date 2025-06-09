local isCaps(chr) = std.codepoint(chr) >= 65 && std.codepoint(chr) <= 90;
local isNum(chr) = std.codepoint(chr) >= 48 && std.codepoint(chr) <= 57;
local Snake(str) = std.lstripChars(
  std.asciiLower(
    std.join(
      '', std.mapWithIndex(
        function(i, chr)
          local chrs = std.stringChars(str);
          if i == 0 || chrs[i - 1] == '_'  // first char, chars after _ = just lowercase
          then chr
          else if isCaps(chr)
                  && !isCaps(chrs[i - 1])
          then '_' + chr
          else if isCaps(chr)
                  && std.length(str) > i + 1
                  && !isCaps(chrs[i + 1])
                  && !isNum(chrs[i + 1])
          then '_' + chr
          else chr
        , std.stringChars(str)
      )
    )
  ), '_'
);

local Object(name, type) = {
  paths+: {
    ['/' + Snake(name) + 's']+: {
      put+: {
        summary: 'Create a %s object' % name,
        description: 'Create a %s object' % name,
        requestBody+: {
          required: true,
          content+: {
            'application/json'+: {
              schema+: {
                '$ref': '#/components/schemas/' + name,
              },
            },
          },
        },
        responses+: {
          '201'+: {
            description: 'Created',
            content+: {
              'application/json'+: {
                schema+: {
                  '$ref': '#/components/schemas/' + name,
                },
              },
            },
            headers: {
              Location: {
                schema: {
                  type: 'string',
                },
                description: 'Location of the created %s object.' % name,
              },
            },
          },
          default: {
            '$ref': '#/components/responses/Error',
          },
        },
      },
      get+: {
        summary: 'List all %s objects' % name,
        description: 'List all %s objects' % name,
        responses+: {
          '200'+: {
            description: 'List of %s objects' % name,
            content+: {
              'application/json'+: {
                schema+: {
                  type: 'array',
                  items+: {
                    '$ref': '#/components/schemas/' + name,
                  },
                },
              },
            },
          },
          default+: {
            '$ref': '#/components/responses/Error',
          },
        },
      },
    },
    ['/' + Snake(name) + 's/{id}']+: {
      put+: {
        summary: 'Update a %s object' % name,
        description: 'Update a %s object' % name,
        parameters+: [
          {
            name: 'id',
            description: 'ID of the %s object' % name,
            'in': 'path',
            required: true,
            schema+: {
              type: 'integer',
              format: 'int64',
            },
          },
        ],
        requestBody+: {
          required: true,
          content+: {
            'application/json'+: {
              schema+: {
                '$ref': '#/components/schemas/' + name + 'Update',
              },
            },
          },
        },
        responses+: {
          '204'+: {
            description: 'Updated',
          },
          '304'+: {
            description: 'Not Modified',
            content+: {
              'application/json'+: {
                schema+: {
                  '$ref': '#/components/schemas/' + name,
                },
              },
            },
          },
          '404'+: {
            '$ref': '#/components/responses/NotFound',
          },
          default: {
            '$ref': '#/components/responses/Error',
          },
        },
      },
      get+: {
        summary: 'Get a %s object' % name,
        description: 'Get a %s object' % name,
        parameters+: [
          {
            name: 'id',
            description: 'ID of the enum',
            'in': 'path',
            required: true,
            schema+: {
              type: 'integer',
              format: 'int64',
            },
          },
        ],
        responses+: {
          '200'+: {
            description: name + ' Object Response',
            content+: {
              'application/json'+: {
                schema+: {
                  '$ref': '#/components/schemas/' + name,
                },
              },
            },
          },
          '404'+: {
            '$ref': '#/components/responses/NotFound',
          },
        },
      },
      delete+: {
        summary: 'Delete a %s object' % name,
        description: 'Get a %s object' % name,
        parameters+: [
          {
            name: 'id',
            description: 'ID of the enum',
            'in': 'path',
            required: true,
            schema+: {
              type: 'integer',
              format: 'int64',
            },
          },
        ],
        responses+: {
          '204'+: {
            description: name + ' Object Deletion Response',
          },
          '404'+: {
            '$ref': '#/components/responses/NotFound',
          },
        },
      },
    },
  },
  components+: {
    schemas+: {
      [name]+: type,
      [name + 'Update']: type { required:: [] },
    },
  },
};
local Common(name) = {
  properties+: {
    id+: {
      type: 'integer',
      format: 'int64',
      readOnly: true,
      description: 'Unique ID of the %s' % name,
    },
    source+: {
      type: 'string',
      description: 'Identity of the user ro tool that created the %s (read-only)' % name,
      readOnly: true,
    },
    name+: {
      type: 'string',
      description: 'User-defined name of the %s' % name,
    },
    client_labels+: {
      type: 'object',
      description: 'Arbitrary key/value metadata about the %s, set by the submitter' % name,
      additionalProperties+: {
        type: 'string',
      },
      default: {},
    },
    server_labels+: {
      type: 'object',
      description: 'Arbitrary key/value metadata about the %s, set by the system' % name,
      additionalProperties+: {
        type: 'string',
      },
      readOnly: true,
      default: {},
    },
  },
  required+: ['id', 'source', 'client_labels', 'server_labels'],
};


{
  openapi: '3.0.0',
  info: {
    title: 'POE Schema Claims API',
    version: '1.0.0',
    description: 'A declarative management API for community schema definition of Path of Exile tabular data files. This API manages `ColumnClaims`, which represent columns in `.datc64` files used in Path of Exile and Path of Exile 2',
    contact: {
      name: 'Graham Forest',
      email: 'graham.forest@protonmail.com',
      url: 'https://poe-schema.obsoleet.org/',
    },
  },
  servers: [
    {
      url: '/v1/',
    },
    {
      url: 'https://poe-schema.obsoleet.org/v1',
      description: 'Production server',
    },
  ],
  paths: {},
  components: {
    schemas: {
      Error: {
        description: 'Generic Error Body',
        type: 'object',
        properties: {
          code: {
            type: 'integer',
          },
          message: {
            type: 'string',
          },
        },
        required: [
          'code',
          'message',
        ],
      },
      NotFound: {
        description: 'Not Found Error Body',
        type: 'object',
        properties: {
          code: {
            type: 'integer',
            description: 'Error code',
          },
          message: {
            type: 'string',
            description: 'Error message',
          },
        },
        required: [
          'code',
          'message',
        ],
      },
      Scalar: {
        oneOf: [
          {
            type: 'object',
            properties: {
              type: {
                type: 'string',
                enum: [
                  'Unknown',
                  'SelfRow',
                  'Bool',
                  'String',
                  'I16',
                  'U16',
                  'I32',
                  'U32',
                  'F32',
                  'I64',
                  'U64',
                ],
              },
            },
            required: [
              'type',
            ],
          },
          {
            type: 'object',
            properties: {
              type: {
                type: 'string',
                enum: [
                  'EnumRow',
                  'ForeignRow',
                  'RowRef',
                ],
              },
              target: {
                type: 'string',
                description: 'The name of the table or enum this column points to',
                example: 'mods',
              },
            },
            required: [
              'type',
              'target',
            ],
          },
        ],
      },
    },
    responses: {
      Error: {
        description: 'Generic Error',
        content: {
          'application/json': {
            schema: {
              '$ref': '#/components/schemas/Error',
            },
          },
        },
      },
      NotFound: {
        description: 'Not Found Error - Used by endpoints that explicitly return a 404 error',
        content: {
          'application/json': {
            schema: {
              '$ref': '#/components/schemas/NotFound',
            },
          },
        },
      },
    },
  },
}
+ Object('Enum', Common('Enum') {
  type: 'object',
  properties+: {
    values+: {
      type: 'array',
      description: 'Ordered list of Enum values',
      items+: {
        type: 'string',
      },
    },
    zero_indexed+: {
      type: 'boolean',
      description: 'Whether the Enum is zero-indexed (true) or one-indexed',
      default: true,
    },
  },
  required+: ['name', 'values', 'zero_indexed'],
})
+ Object('ColumnClaim', Common('ColumnClaim') {
  type: 'object',
  properties+: {
    offset+: {
      type: 'integer',
      description: 'Byte offset of the start of the ColumnClaim in each row',
      minimum: 0,
    },
    bytes+: {
      type: 'integer',
      description: 'Number of bytes used by the ColumnClaim, starting at the offset',
      minimum: 1,
      maximum: 8,
    },
    is_array+: {
      type: 'boolean',
      description: 'Whether the ColumnClaim is an array',
      default: false,
    },
    column+: {
      '$ref': '#/components/schemas/Scalar',
    },
    datfile+: {
      type: 'string',
      description: 'The datfile basename that this ColumnClaim references',
      example: 'mods',
    },
  },
  required+: ['offset', 'bytes', 'column', 'datfile'],
})

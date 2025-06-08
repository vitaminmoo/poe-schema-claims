local api = {
  Object(name, type):: {
    paths+: {
      ['/' + std.asciiLower(name) + 's']+: {
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
      ['/' + std.asciiLower(name) + 's' + '/{id}']+: {
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
                type: 'string',
              },
            },
          ],
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
            '200'+: {
              description: 'Updated',
              content+: {
                'application/json'+: {
                  schema+: {
                    '$ref': '#/components/schemas/' + name,
                  },
                },
              },
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
                type: 'string',
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
                type: 'string',
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
      },
    },
  },
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
      Cell: {
        type: 'object',
        properties: {
          type: {
            '$ref': '#/components/schemas/Scalar',
          },
          is_array: {
            type: 'boolean',
            description: 'Whether the cell is an array',
          },
        },
        required: ['type'],
      },
      Scalar: {
        oneOf: [
          {
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
          {
            type: 'object',
            properties: {
              type: {
                type: 'string',
                enum: [
                  'ForeignRow',
                ],
              },
              target: {
                type: 'string',
                description: 'The target datfile this foreign row points to',
              },
            },
            required: [
              'type',
              'target',
            ],
          },
          {
            type: 'object',
            properties: {
              type: {
                type: 'string',
                enum: [
                  'EnumRow',
                ],
              },
              enum_name: {
                type: 'string',
                description: 'The name of the enum this row points to',
              },
            },
            required: [
              'type',
              'enum_name',
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
+ api.Object('Enum', {
  oneOf: [
    {
      type: 'object',
      properties: {
        id: {
          type: 'string',
          readOnly: true,
          description: 'Unique ID of the enum',
        },
        values: {
          type: 'array',
          description: 'Ordered list of enum values',
          items: {
            type: 'string',
          },
        },
        zero_indexed: {
          type: 'boolean',
          description: 'Whether the enum is zero-indexed (true) or one-indexed',
          default: true,
        },
      },
      required: ['id', 'values', 'zero_indexed'],
    },
  ],
})
+ api.Object('ColumnClaim', {
  oneOf: [
    {
      type: 'object',
      properties: {
        id: {
          type: 'string',
          readOnly: true,
          description: 'Unique ID of the column claim',
        },
        name: {
          type: 'string',
          description: 'User-defined name of the column claim (optional)',
        },
        offset: {
          type: 'integer',
          description: 'Byte offset of the start of the ColumnClaim in each row',
          minimum: 0,
        },
        bytes: {
          type: 'integer',
          description: 'Number of bytes used by the ColumnClaim, starting at the offset',
          minimum: 1,
          maximum: 8,
        },
        labels: {
          type: 'object',
          description: 'Arbitrary key/value metadata about the ColumnClaim',
          additionalProperties: {
            type: 'string',
          },
        },
        column_type: {
          '$ref': '#/components/schemas/Cell',
        },
        source: {
          type: 'string',
          description: 'Identity of the user ro tool that created the ColumnClaim (read-only)',
          readOnly: true,
        },
        datfile: {
          type: 'string',
          description: 'The datfile basename that this ColumnClaim references',
          example: 'mods',
        },
      },
      required: ['offset', 'bytes', 'column_type', 'source', 'datfile'],
    },
  ],
})

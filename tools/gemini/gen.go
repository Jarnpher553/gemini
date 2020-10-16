package main

const gen = `{
  "name": "default",
  "path": "./",
  "services": [
    {
      "name": "device",
      "orm": {
        "name": "device",
        "fields": [
          {
            "name": "id",
            "type": "uuid",
            "nullable": false
          },
          {
            "name": "code",
            "type": "string",
            "nullable": false
          },
          {
            "name": "serial",
            "type": "string",
            "nullable": true
          },
          {
            "name": "from",
            "type": "time",
            "nullable": true
          }
        ]
      },
      "dto": {
        "request": {
          "name": "deviceIn",
          "fields": [
            {
              "name": "code",
              "type": "string",
              "required": true
            },
            {
              "name": "serial",
              "type": "string",
              "required": true
            },
            {
              "name": "from",
              "type": "time",
              "required": true
            }
          ]
        },
        "response": {
          "name": "deviceOut",
          "fields": [
            {
              "name": "id",
              "type": "uuid"
            },
            {
              "name": "code",
              "type": "string"
            },
            {
              "name": "serial",
              "type": "string"
            }
          ]
        }
      }
    }
  ]
}
`

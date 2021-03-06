# GENERATED BY KOMAND SDK - DO NOT EDIT
import komand
import json


class EmitGreetingInput(komand.Input):
    schema = json.loads("""
   {
  "type": "object",
  "title": "Variables",
  "properties": {
    "frequency": {
      "type": "integer",
      "title": "Frequency",
      "description": "How frequently (in seconds) to trigger a greeting",
      "default": 15,
      "order": 1
    }
  }
}
    """)

    def __init__(self):
        super(self.__class__, self).__init__(self.schema)


class EmitGreetingOutput(komand.Output):
    schema = json.loads("""
   {
  "type": "object",
  "title": "Variables",
  "properties": {
    "greeting": {
      "type": "string",
      "title": "Greeting",
      "order": 1
    },
    "time": {
      "type": "string",
      "title": "Time",
      "displayType": "date",
      "format": "date-time",
      "order": 2
    }
  }
}
    """)

    def __init__(self):
        super(self.__class__, self).__init__(self.schema)

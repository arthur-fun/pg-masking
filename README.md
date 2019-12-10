# pg-masking
This is a tool to masking data of postgres database due to security reason. For now, only postgresql is supported.

# usage
$ pg_masking -f config.json

# configuration
The tool pg_masking needs a configuration file in json format, as shown below:
```json5
{
    "tables": [
        "table1", "table2", "table3"
    ],
    "source": {
        "connection": {
            "dbms": "postgres",
            "dburl": ""
        }
    },
    "column-converter": [
        {
            "table-name": "*",
            "column-name": "column1",
            "converter": "converter1"
        },
        {
            "table-name": "table2",
            "column-name": "column2",
            "converter": "converter2"
        }
    ],
    "destination": {
        "connection": {
            "dbms": "postgres",
            "dburl": ""
        }
    }
}
```


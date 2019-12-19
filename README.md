# pg-masking
This is a tool to masking data of postgres database due to security reason. For now, only postgresql is supported.

# usage
$ pg_masking -f config.json

# configuration
The tool pg_masking needs a configuration file (e.g. config.json) in json format, as shown below:
```json5
{
    "tables": [
        "t1"
    ],
    "source": {
        "dbms": "postgres",
        "dburl": "postgres://<user>:<password>@<host>/<db1>?sslmode=disable"
    },
    "destination": {
        "dbms": "postgres",
        "dburl": "postgres://<user>:<password>@<host>/<db2>?sslmode=disable"
    },
    "column-converter": [
        {
            "table-name": "*",
            "column-name": "amount",
            "converter": "Random"
        },        
        {
            "table-name": "*",
            "column-name": "password",
            "converter": "ReplaceAll",
            "converter-parameters": ""
        },
        {
            "table-name": "t1",
            "column-name": "phone_number",
            "converter": "Replace",
            "converter-parameters": "*****, 4"
        }
    ]
}
```

# database data type supported to be transfered for now is shown below:
- "BOOL"          
- "BPCHAR"        
- "CHAR"          
- "DATE"          
- "FLOAT4"        
- "FLOAT8"        
- "INT2"          
- "INT4"          
- "INT8"          
- "NUMERIC"       
- "TEXT"          
- "TIME"          
- "TIMETZ"        
- "TIMESTAMP"     
- "TIMESTAMPTZ"   
- "VARCHAR"       

# converters and database data type supported to be masked for now is shown below:
- ReplaceAll: BPCHAR, CHAR, VARCHAR, TEXT
- Replace: BPCHAR, CHAR, VARCHAR, TEXT
- Random: INT2, INT4, INT8, "NUMERIC", "FLOAT4", "FLOAT8"

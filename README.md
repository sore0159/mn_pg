### Manage PostGreSQL

``mn_pg`` is a commandline utility to simplify the creation and deletion of tables in postgresql.  More sophisticated ALTER TABLE schema management is not within the scope of this project.

### JSON Configuration
You must provide a json file with the database information mn\_pg will need to connect to postgresql, and the location of the sql files mn_pg will use for the table schemas.  

An example.json file is included in the repository.  The json file must be of the following form:
```json
{
        "dbName": "example_db",
        "dbUser": "example_user",
        "dbPassword": "example_pass",
        "schemaFilenames": ["example.sql"],

        "defaultDropAll": false,
        "defaultDropSome": ["exampletable"],

        "defaultCreateAll": false,
        "defaultCreateSome": ["exampletable"]
}
```

Any value with the prefix 'default' is optional and will control mn\_pg's behavior when run without command line flags.  Options to drop all or create all will use the supplied sql files to determine which tables to drop and create.

### Command line flags
Alternatively, tables to drop or create may be specified via commandline flags.  These flags will override defaults present in the json file.

```bash
mn_pg -d [TABLENAMES...] -c [TABLENAMES...]
mn_pg -dc
```
The -d flag will specify tables to drop; if none are listed following the flag then all tables in the sql files will be dropped.  Likewise the -c flag will specify tables to create (their definitions must be found in the provided sql files) and default to all if none are listed on the commandline.  Both -d and -c flags can be present multiple times, and can appear in any order.  

A sole flag of -dc will set all tables to be dropped and created.

A flag of -h anywhere will cause mn\_pg to do nothing but print a help message.  Any error in parsing user input will cause mn\_pg to display the error, a help message, and then abort.


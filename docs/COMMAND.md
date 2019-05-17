#### ENGINE COMMANDS

| COMMAND        | DESCRIPTION                          |
|----------------|--------------------------------------|
| info session   | Print current session information    |
| info api       | Print api rest endpoints information |
| env            | Print environment variable           |
| help           | Print help information               |
| clear          | Clear current screen                 |
| api <run/stop> | (Run/Stop) restful API               |


#### TARGET COMMANDS

| COMMAND                                       | DESCRIPTION                                 |
|-----------------------------------------------|---------------------------------------------|
| target add <type> <value>                     | Add new target                              |
| target view result <target_id> <result_id>    | View one result from targets                |
| target view results <target_id> <module_name> | View all targets results from module name   |
| target links <target_id>                      | View linked targets                         |
| target update <target_id> <value>             | Update a target                             |
| target delete <target_id>                     | Remove target by ID                         |
| target list                                   | List subjects                               |
| target modules <target_id>                    | List modules available with selected target |

#### NOTE COMMANDS

| COMMAND                            | DESCRIPTION                           |
|------------------------------------|---------------------------------------|
| note add <id target/result> <text> | Add new note to target or result      |
| note view <id target/result>       | View note linked to target or result  |

#### MODULE COMMANDS
| COMMAND                         | DESCRIPTION           |
|---------------------------------|-----------------------|
| <module> target <target_id>     | Set a target argument |
| <module> filter <filter>        | Set a filter argument |
| <module> set <argument> <value> | Set specific argument |
| <module> list                   | List module arguments |
| <module> run                    | Run selected module   |

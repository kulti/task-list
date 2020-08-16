# Creating a new migration

## Generate migrations' names

```
./script/deploy.sh dev gen-migration-name test_table
```

## Write migration rules

It is hard to support down migrations. It is harder to support applying down migration.

So the simplest solution is: NEVER do an up migration that breaks backward compatibility. It means always add tables and columns and never drop or rename.

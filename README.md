# tfmig

tfmig is the tool to migrate terraform state between project.  
you can move one or more number of states that you want to move at once.  
you can select states by fazzyfinder CUI interface.  

# usage

```
% tfmig -s /path/to/move-from -d /path/to/move-to [-w workspace-name]
```

- s: path to terraform project that you wanna move state from
- d: path to terraform project that you wanna move state into
- w (optional): workspace name

# licence

MIT

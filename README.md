# tfmig

tfmig is the tool to migrate terraform state between project.  
you can move only one state at once per execution.  
you can choose one state in fazzyfinder CUI interface.  

# usage

```
% tfmig -s /path/to/move-from -d /path/to/move-to (-w sandbox)
```

- s: path to terraform project that you wanna move state from
- d: path to terraform project that you wanna move state into
- w (optional): workspace for use

# licence

MIT

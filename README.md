# Launch or Focus

Simple utility tool for [Niri](https://github.com/YaLTeR/niri) that allows to focus a window by App ID or run a command if App ID is not found.

## Examples

### Alacritty 

```bash
lof Alacritty alacritty
```

if any window with app ID "Alacritty" is found, it will be focused. Otherwise the command will execute "alacritty"

### Thunderbird binding

```kdl
binds {
  Mod+E repeat=false {spawn-sh "lof org.mozilla.Thunderbird thunderbird"; }
}
```

Binding Thunderbird to Mod+e with lof means that with this binding I will always have only one instance of thunderbird running, and I don't need to worry where I spawned it initially because it will focus automatically

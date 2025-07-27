local g = import '../g.libsonnet';

{
  stat: {
    local stat = g.panel.stat,

    small(title, description):
      stat.new(title)
      + stat.panelOptions.withDescription(description)
      + stat.panelOptions.withGridPos("3", "4"),

    large(title, description):
      stat.new(title)
      + stat.panelOptions.withDescription(description)
      + stat.panelOptions.withGridPos("6", "4"),
  }
}

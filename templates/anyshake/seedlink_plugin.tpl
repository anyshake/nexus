template: $template

plugin $seedlink.source.id cmd="$seedlink.plugin_dir/anyshake_plugin -address $sources.anyshake.address -timeout $sources.anyshake.timeout"
             timeout = 600
             start_retry = 60
             shutdown_wait = 10

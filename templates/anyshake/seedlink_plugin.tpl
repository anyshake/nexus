template: $template

plugin $seedlink.source.id cmd="$seedlink.plugin_dir/anyshake_plugin -address $sources.anyshake.address -timeout $sources.anyshake.timeout -realtime=$sources.anyshake.realtime -verbose $sources.anyshake.verbose"
             timeout = 5
             start_retry = 30
             shutdown_wait = 10

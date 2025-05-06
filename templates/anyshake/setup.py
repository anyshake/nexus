'''
Plugin handler for the AnyShake plugin.
'''


class SeedlinkPluginHandler:
    def __init__(self): pass

    def push(self, seedlink):
        # Set default address
        try:
            seedlink.param('sources.anyshake.address')
        except:
            seedlink.setParam('sources.anyshake.address', '127.0.0.1:30000')

        # Set default timeout
        try:
            seedlink.param('sources.anyshake.timeout')
        except:
            seedlink.setParam('sources.anyshake.timeout', 10)

        # Use network.station as unique key
        return seedlink.net + "." + seedlink.sta

    def flush(self, seedlink):
        pass

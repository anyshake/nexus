'''
Plugin handler for the AnyShake plugin.
'''


class SeedlinkPluginHandler:
    def __init__(self): pass

    def push(self, seedlink):
        try:
            address = seedlink.param('sources.anyshake.address')
        except:
            address = '127.0.0.1:30000'
            seedlink.setParam('sources.anyshake.address', address)

        try:
            seedlink.param('sources.anyshake.timeout')
        except:
            seedlink.setParam('sources.anyshake.timeout', 10)

        try:
            realtime = seedlink.param('sources.anyshake.realtime')
        except:
            realtime = 'false'
            seedlink.setParam('sources.anyshake.realtime', realtime)

        try:
            seedlink.param('sources.anyshake.verbose')
        except:
            seedlink.setParam('sources.anyshake.verbose', "false")

        return (address, realtime)

    def flush(self, seedlink):
        pass

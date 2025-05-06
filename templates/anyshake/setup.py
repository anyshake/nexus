'''
Plugin handler for the AnyShake plugin.
'''


class SeedlinkPluginHandler:
    def __init__(self): pass

    def push(self, seedlink):
        try:
            seedlink.param('sources.anyshake.address')
        except:
            seedlink.setParam('sources.anyshake.address', '127.0.0.1:30000')

        try:
            seedlink.param('sources.anyshake.timeout')
        except:
            seedlink.setParam('sources.anyshake.timeout', 10)

        try:
            seedlink.param('sources.anyshake.verbose')
        except:
            seedlink.setParam('sources.anyshake.verbose', "false")

        try:
            seedlink.param("sources.anyshake.proc")
        except:
            seedlink.setParam("sources.anyshake.proc", "e_c1x1g_6ch_250hz")

        return None

    def flush(self, seedlink):
        pass

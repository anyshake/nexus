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
            seedlink.param('sources.anyshake.sampleRate')
        except:
            seedlink.setParam('sources.anyshake.sampleRate', 250)

        try:
            seedlink.param('sources.anyshake.locationCode')
        except:
            seedlink.setParam('sources.anyshake.locationCode', '00')

        try:
            seedlink.param('sources.anyshake.channelPrefixVelocity')
        except:
            seedlink.setParam('sources.anyshake.channelPrefixVelocity', 'EH')

        try:
            seedlink.param('sources.anyshake.channelPrefixAcceleration')
        except:
            seedlink.setParam(
                'sources.anyshake.channelPrefixAcceleration', 'EN')

        try:
            seedlink.param('sources.anyshake.verbose')
        except:
            seedlink.setParam('sources.anyshake.verbose', "false")

        try:
            seedlink.param("sources.anyshake.proc")
        except:
            seedlink.setParam("sources.anyshake.proc", "e_c1x1g_6ch")

        return None

    def flush(self, seedlink):
        pass

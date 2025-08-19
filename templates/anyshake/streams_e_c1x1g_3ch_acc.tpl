<proc name="e_c1x1g_3ch_acc">
    <tree>
        <input name="CH4" channel="Z" location="$sources.anyshake.locationCode" rate="$sources.anyshake.sampleRate" />
        <input name="CH5" channel="E" location="$sources.anyshake.locationCode" rate="$sources.anyshake.sampleRate" />
        <input name="CH6" channel="N" location="$sources.anyshake.locationCode" rate="$sources.anyshake.sampleRate" />
        <node stream="$sources.anyshake.channelPrefixAcceleration" />
    </tree>
</proc>

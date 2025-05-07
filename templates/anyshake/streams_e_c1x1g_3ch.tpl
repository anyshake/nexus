<proc name="e_c1x1g_3ch">
    <tree>
        <input name="CH1" channel="Z" location="$sources.anyshake.locationCode" rate="$sources.anyshake.sampleRate" />
        <input name="CH2" channel="E" location="$sources.anyshake.locationCode" rate="$sources.anyshake.sampleRate" />
        <input name="CH3" channel="N" location="$sources.anyshake.locationCode" rate="$sources.anyshake.sampleRate" />
        <node stream="EH" />
    </tree>
</proc>

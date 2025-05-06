<proc name="e_c1x1g_6ch_200hz">
    <tree>
        <input name="EHZ" channel="Z" location="" rate="200" />
        <input name="EHE" channel="E" location="" rate="200" />
        <input name="EHN" channel="N" location="" rate="200" />
        <node stream="EH" />
    </tree>
    <tree>
        <input name="ENZ" channel="Z" location="" rate="200" />
        <input name="ENE" channel="E" location="" rate="200" />
        <input name="ENN" channel="N" location="" rate="200" />
        <node stream="EN" />
    </tree>
</proc>

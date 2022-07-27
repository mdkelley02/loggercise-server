<h2>Loggercise Server</h2>
<h5>.env</h5>
<pre>
MONGO_CONN_STRING='mongodb+srv://USERNAME:PASSOWRD!@URL/?retryWrites=true&w=majority'
</pre>
<h5>Run</h5>
<pre>
cd app/ && make run
</pre>
<h5>Protos</h5>
<pre>
cd app/protos && buf generate
</pre>
<h5>Deploy</h5>
<pre>
terraform -chdir=build init
terraform -chdir=build apply
terraform -chdir=build destroy
</pre>
{{define "index"}}

<section class="blocks">

	<div class="row">
		<div class="col-lg-3 offset-lg-1 visible-lg">
			<div class="list-group">
				{{range $i, $mod := .user.Modules}}
				<a href="{{$mod.Module_path}}" class="list-group-item">{{$mod.Module}}</a>
				{{end}}
			</div>
		</div>
		<div class="col-lg-8">
			<h4>CURT Data Caching</h4>
			<div class="well redis-search">
				<form class="form-inline" role="form">
					<div class="row">
						<div class="form-group col-xs-4">
							<label for="redis-namespace" class="sr-only">Select Namespace</label>
							<select name="redis-namespace" id="redis-namespace" class="form-control">
								<option value="">- Select Namespace -</option>
								{{range $i, $nps := .RedisNamespaces}}
								<option value="{{$i}}">{{$i}}</option>
								{{end}}
							</select>
						</div>
						<div class="form-group col-xs-6">
							<label for="redis-key" class="sr-only">Find Key</label>
							<input type="search" class="input-sm form-control redis-key" id="redis-key" placeholder="Search for key..">
						</div>
						<div class="form-group col-xs-1">
							<input type="submit" class="btn btn-primary" value="Search">
						</div>
					</div>
				</form>
			</div>

			<div class="table-responsive">
				<table class="table sortable cache-table">
					<thead>
						<tr>
							<th>Key</th>
							<th>Value</th>
							<th>Expires (seconds)</th>
							<th data-defaultsort="disabled"></th>
						</tr>
					</thead>
					<tbody>
						{{range $i, $data := .RedisData}}
						<tr>
							<td>{{$data.Key}}</td>
							<td>{{$data.Value}}</td>
							<td>{{$data.Ttl}}</td>
							<td>
								<div class="btn-group">
									<button class="btn btn-primary dropdown-toggle" type="button" data-toggle="dropdown">
										Action
										<span class="caret"></span>
									</button>
									<ul class="dropdown-menu" role="menu">
										<li>
											<a href="#" data-key="{{$data.Key}}" class="redis-preview">View</a>
										</li>
										<li>
											<a href="#" data-key="{{$data.Key}}" class="redis-delete">Delete</a>
										</li>
									</ul>
								</div>
							</td>
						</tr>
						{{end}}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</section>

{{end}}
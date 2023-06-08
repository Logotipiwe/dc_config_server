package pages

import "net/http"

func GetAdminPage(w http.ResponseWriter, r *http.Request) string {
	return `
<h2>THIS IS CONFIG SERVER!!</h2>
<div id="service-select-form">
<select>
	<option value="*">all</option>
	<option value="service-one">service-one</option>
	<option value="service-two">service-two</option>
</select>
<input type="submit" value="choose service">
<h2>Properties:</h2>
<div class="property">
	<input value="SOME_PROP">
	<input value="value of this prop">
	<form action="/delete_prop" method="post">
		<input type="hidden" name="id" value="364576572967386">
		<input type="submit" value="Delete">
	</form>
</div>
</div>
<style>
.property{
	display: flex;
}
</style>
`
}

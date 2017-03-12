<form id="seq">
<label for="taxon">Taxon: </label><input type="text" name="taxon">

<label for="search_term">Search Term: </label><input id="search_term" name="search_term" type="text" list="search_terms" />
<datalist id="search_terms">
   <option value="COI">
   <option value="&quot;(COI[All Fields] OR \&quot;cytochrome oxidase I\&quot;[All Fields] OR \&quot;cytochrome oxidase subunit I\&quot;[All Fields] OR COX1[All Fields] OR \&quot;COXI\&quot;[All Fields]) NOT (\&quot;complete genome\&quot;[title] OR \&quot;complete DNA\&quot;[title])&quot;">
</datalist>
<select name="retmax">
  <option value="10" selected>10</option>
  <option value="50">50</option>
  <option value="100">100</option>
  <option value="500">500</option>
</select>
  <input type="submit" value="Submit">
  <input type="button" onclick="myFunction()" value="Reset form">
</form>
<script>
function myFunction() {
    document.getElementById("seq").reset();
}
</script>
Submitted:
<br>
<br>
<?php
	
	if(isset($_GET['search_term'])) {
		echo $_GET["taxon"];
		echo "<br>";
		echo $_GET["search_term"];
		echo "<br>";
		$command = "sequence-manager.exe -sterm=$_GET[search_term] -taxon=$_GET[taxon] -retmax=$_GET[retmax]";
		echo $command;
		echo "<br>";
		echo "<br>";
		$command_string = "sequence-manager.exe";
		if (strtoupper(substr(PHP_OS, 0, 3)) === 'WIN') {
		} else {
			$command_string = "./sequence-manager";
		}
		system("$command_string -sterm=$_GET[search_term] -taxon=$_GET[taxon] -retmax=$_GET[retmax]");
		
		header( 'Location: ../sequence-manager/index.php' ) ;
	}
	
	$dir    = '../sequence-manager/output';
	$files = scandir($dir);
	$files = array_reverse($files);
	
	foreach ($files as &$file) {
		if (strpos($file, '.html') !== false) {
			echo "<a href=\"output\\" .$file . "\">" . $file . "</a> |" ;
		}
		if (strpos($file, '.csv') !== false) {
			echo "<a href=\"output\\" .$file . "\">" . $file . "</a> <br>" ;
		}
		
	}

?>
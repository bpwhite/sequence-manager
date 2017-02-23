<form>
<label for="taxon">Taxon: </label><input type="text" name="taxon">

<label for="search_term">Search Term: </label><input id="search_term" name="search_term" type="text" list="search_terms" />
<datalist id="search_terms">
   <option value="COI">
   <option value="(COI[All Fields] OR &quot;cytochrome oxidase I&quot;[All Fields] OR &quot;cytochrome oxidase subunit I&quot;[All Fields] OR COX1[All Fields] OR &quot;COXI&quot;[All Fields]) NOT (&quot;complete genome&quot;[title] OR &quot;complete DNA&quot;[title])">
</datalist>
  <input type="submit" value="Submit">
</form>

Submit: 
<?php
	echo $_GET["taxon"];
	echo "<br>";
	echo $_GET["search_term"];
	system("sequence-manager.exe -sterm=$_GET[search_term] -taxon=$_GET[taxon] -retmax=20");
?>
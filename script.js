function generateDungeon() {
  d = new Date();
  $("#dungeonImage").attr("src", "/generate?"+d.getTime());
}

$("#generateButton").click(generateDungeon);
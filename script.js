var dungeon_width = 100;
$(function() {
	$( "#dungeon_width_slider" ).slider({
		value: dungeon_width,
		min: 20,
		max: 1000,
		step: 10,
		slide: function( event, ui ) {
			$( "#dungeon_width_label" ).html( "Dungeon width: " + ui.value );
		}
	});
});
$( "#dungeon_width_label" ).html( "Dungeon width: " + dungeon_width );

var dungeon_height = 100;
$(function() {
	$( "#dungeon_height_slider" ).slider({
		value: dungeon_height,
		min: 20,
		max: 1000,
		step: 10,
		slide: function( event, ui ) {
			$( "#dungeon_height_label" ).html( "Dungeon height: " + ui.value );
		}
	});
});
$( "#dungeon_height_label" ).html( "Dungeon height: " + dungeon_height );

var room_attempts = 200;
$(function() {
	$( "#room_attempts_slider" ).slider({
		value: room_attempts,
		min: 10,
		max: 1000,
		step: 10,
		slide: function( event, ui ) {
			$( "#room_attempts_label" ).html( "Room attempts: " + ui.value );
		}
	});
});
$( "#room_attempts_label" ).html( "Room attempts: " + room_attempts );

var min_room_size = 5;
$(function() {
	$( "#min_room_size_slider" ).slider({
		value: min_room_size,
		min: 1,
		max: 100,
		step: 1,
		slide: function( event, ui ) {
			$( "#min_room_size_label" ).html( "Min room size: " + ui.value );
		}
	});
});
$( "#min_room_size_label" ).html( "Min room size: " + min_room_size );

var max_room_size = 15;
$(function() {
	$( "#max_room_size_slider" ).slider({
		value: max_room_size,
		min: 1,
		max: 100,
		step: 1,
		slide: function( event, ui ) {
			$( "#max_room_size_label" ).html( "Max room size: " + ui.value );
		}
	});
});
$( "#max_room_size_label" ).html( "Max room size: " + max_room_size );

var pixel_size = 4;
$(function() {
	$( "#pixel_size_slider" ).slider({
		value: pixel_size,
		min: 1,
		max: 20,
		step: 1,
		slide: function( event, ui ) {
			$( "#pixel_size_label" ).html( "Pixel size: " + ui.value );
		}
	});
});
$( "#pixel_size_label" ).html( "Pixel size: " + pixel_size );

function generateDungeon() {
	d = new Date();
	queryString = "/generate/?time="+d.getTime();

	dungeonWidth = $( "#dungeon_width_slider" ).slider( "option", "value" );
	if(dungeonWidth) {
		queryString += "&dungeonWidth=" + dungeonWidth;
	}

	dungeonHeight = $( "#dungeon_height_slider" ).slider( "option", "value" );
	if(dungeonHeight) {
		queryString += "&dungeonHeight=" + dungeonHeight;
	}

	roomAttempts = $( "#room_attempts_slider" ).slider( "option", "value" );
	if(roomAttempts) {
		queryString += "&roomAttempts=" + roomAttempts;
	}

	minRoomSize = $( "#min_room_size_slider" ).slider( "option", "value" );
	if(minRoomSize) {
		queryString += "&minRoomSize=" + minRoomSize;
	}

	maxRoomSize = $( "#max_room_size_slider" ).slider( "option", "value" );
	if(maxRoomSize) {
		queryString += "&maxRoomSize=" + maxRoomSize;
	}

	pixelSize = $( "#pixel_size_slider" ).slider( "option", "value" );
	if(pixelSize) {
		queryString += "&pixelSize=" + pixelSize;
	}

	try {
		seed = parseInt( $( "#seed" ).val() );
		if(!isNaN(seed) && seed % 1 === 0) {
			queryString += "&seed=" + seed
		} else {
			console.log("Invalid seed!");
		}
	} catch(e) {
		// invalid number
		console.log("Invalid seed!");
		console.log(e);
	}

	$("#dungeonImage").attr("src", queryString);
}

$("#generateButton").click(generateDungeon);

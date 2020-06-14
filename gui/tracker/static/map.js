// Getting URL params
const urlParams = new URLSearchParams(window.location.search);
var url = ""

// Generate URL
if (urlParams.has('api') && urlParams.get('port')) {
    url = "http://" + urlParams.get("api") + ":" + urlParams.get('port')
}

// Map layers
var tiles = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 18,
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors, Points &copy 2012 LINZ'
    });

// Default values
latlng = L.latLng(43.604652, 1.444209);
var map = L.map('map', {center: latlng, zoom: 5, layers: [tiles]});
var line = L.polyline([]).addTo(map);
var markerStart = null;
var markerStop = null;
var markersList = null;
var markers = []

// Get Tracked users
getTrackers = function(markers) {
    $.get(url + '/api/trackers', function (data) {
       let trackers = data.data;
       if (trackers !== null) {

           $('#profiles').html('')
           $.each(trackers, function (index, tracker) {
               // Generate a title & icon
               var title = tracker.identifier;
               var icon = L.icon({
                   iconUrl: tracker.picture,
                   iconSize: [45,45],
                   iconAnchor: [22, 94],
                   popupAnchor: [-3, -76],
                   shadowSize: [68, 95],
                   shadowAnchor: [22, 94]
               });

               // Remove marker if already exist
               if (markers[tracker.identifier] !== undefined)
                   markers[tracker.identifier].remove(map)

               var marker = L.marker(new L.LatLng((tracker.position.latitude + (index / 0.001)), tracker.position.longitude), { title: title, icon: icon });
               marker.bindPopup(title);
               marker.addTo(map)

               markers[tracker.identifier] = marker

               let cssClass = 'picture';
               if (tracker.memories.length > 0) {
                   cssClass = cssClass + ' hasMemories'
               }
               let template = '<img class="'+cssClass+'" data-id="'+tracker.id+'" src="'+tracker.picture+'">';

               $('#profiles').append(template);

           });

           markersList = markers;
           map.addLayer(markers);
           clickPicture();
       }
    });
}

function MemoriesDraw(item) {

    line.remove(map)
    line = L.polyline([]).addTo(map);
    if (item.memories.length > 0) {
        $.each($(item.memories), function (i, v) {
            if (i === 0) {
                if (markerStart !== null)
                    markerStart.remove(map)
                markerStart = L.marker(new L.LatLng((v.latitude), v.longitude), {});
                //markerStart.bindPopup(title);
                markerStart.addTo(map)
            }
            point = {lat: v.latitude, lng: v.longitude};
            line.addLatLng(point);
        })
    }

    point = {lat: item.position.latitude, lng: item.position.longitude};
    if (markerStop !== null)
        markerStop.remove(map)
    markerStop = L.marker(new L.LatLng((item.position.latitude), item.position.longitude), {});
    //markerStart.bindPopup(title);
    markerStop.addTo(map)
    line.addLatLng(point);
}

clickPicture = function() {
    $.each($('.picture'), function (index, value) {
        if($(value).attr("onClick") === undefined) {
            $(value).click(function () {
                let id = $(this).attr('data-id')
                let forgedUrl = url + '/api/tracker/' + id
                $.get(forgedUrl, function (data) {
                    if (data.data !== undefined) {
                        MemoriesDraw(data.data)
                    }
                })
            });
        }
    });
}

$(document).ready(function () {

    var markers = L.markerClusterGroup();

    getTrackers(markers);

    setInterval(function () {
        if (markersList !== null)
            markersList.clearLayers();
        getTrackers(markers)
    }, 10000);
})

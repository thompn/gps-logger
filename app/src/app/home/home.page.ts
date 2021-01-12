import { Component, resolveForwardRef } from '@angular/core';
import * as Leaflet from 'leaflet';
import { ApiService } from '../services/api.service';
import { AlertService } from '../services/alert.service';

import { Geolocation } from '@ionic-native/geolocation/ngx';
import { Subscription } from 'rxjs';



@Component({
  selector: 'app-home',
  templateUrl: 'home.page.html',
  styleUrls: ['home.page.scss']
})

export class HomePage {
  map: Leaflet.Map;
  marker: any;
  subscription: Subscription;
  coords: any = [];
  surface: any = "unknown";
  latitude : any;
  longitude : any;

  constructor(
    private geolocation: Geolocation,
    private api: ApiService,
    private alert: AlertService
  ) { }

  private createMap(): void {
    this.map = Leaflet.map('mapId').setView([52.2, 5.4], 8);
    Leaflet.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: 'GPS Trail Logger',
    }).addTo(this.map);
  }

  ionViewDidEnter(): void {
    this.createMap();
    console.log("loaded...")
  }

  ionViewDidLeave(): void {
    this.map.remove();
  }

  btnUpdateLocation(tracking){
    if (tracking){
      tracking = "start"
      this.subscription = this.geolocation.watchPosition().subscribe(async(response: any)=>{
        this.coords = [response.coords.latitude, response.coords.longitude];
        this.api.updateDriverLocation(this.coords[0], this.coords[1], this.surface, tracking)
        .then(()=>{
          if (this.marker != null){
            this.map.removeLayer(this.marker);
          }
          this.marker = Leaflet.marker(this.coords)
          .addTo(this.map).bindPopup('You are here.');
          this.map.setView(this.coords, 16);
        }).catch((error)=>{
          console.log(error);
          this.alert.presentCustomAlert("Error", "Check your internet connection.");
        });
      });
    }else{
      tracking = "stop"
      this.subscription.unsubscribe();
      this.api.updateDriverLocation(this.coords[0], this.coords[1], this.surface, tracking)
      .then(()=>{
        this.alert.presentCustomAlert("Success", "Location paused.");
      }).catch((error)=>{
        console.log(error)
        this.alert.presentCustomAlert("Error", "Check your internet connection.");
      });
    }
  }

btnPaved(surface){
  this.surface = "paved"
  this.btnUpdateLocation(true)

  }

btnDirt(surface){
  this.surface = "dirt"
  this.btnUpdateLocation(true)

}

btnSand(surface){
  this.surface = "sand"
  this.btnUpdateLocation(true)
  }

btnMud(surface){
  this.surface = "mud"
  this.btnUpdateLocation(true)

}
}

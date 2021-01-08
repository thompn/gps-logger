import { Injectable } from '@angular/core';
import { AlertController } from '@ionic/angular';

@Injectable({
  providedIn: 'root'
})
export class AlertService {

  constructor(
    private alert: AlertController
  ) { }

  async presentCustomAlert(title, message): Promise<any>{
    let alert = await this.alert.create({
      header: title,
      message: message,
      buttons: ["OK"]
    });
    alert.present();
  }
}

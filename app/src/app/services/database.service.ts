import { Platform } from '@ionic/angular';
import { Injectable } from '@angular/core';
import { SQLitePorter } from '@ionic-native/sqlite-porter/ngx';
import { HttpClient } from '@angular/common/http';
import { SQLite, SQLiteObject } from '@ionic-native/sqlite/ngx';
import { BehaviorSubject, Observable } from 'rxjs';
 
export interface Dev {
  timestamp: string,
  longitude: string,
  latitude: string,
  surface: string,
  tracking: string
}
 
@Injectable({
  providedIn: 'root'
})
export class DatabaseService {
  private database: SQLiteObject;
  private dbReady: BehaviorSubject<boolean> = new BehaviorSubject(false);
 
  developers = new BehaviorSubject([]);
 
  constructor(
    private plt: Platform, 
    private sqlitePorter: SQLitePorter, 
    private sqlite: SQLite, 
    private http: HttpClient) {
      this.plt.ready().then(() => {
        console.log("ready to seed");
        this.sqlite.create({
          name: 'appdb.db',
          location: 'default'
        }).catch(e => console.error("cannot create",e))
        .then((db: SQLiteObject) => {
            this.database = db;
            this.seedDatabase(); 
        });
    });
    }
 
  seedDatabase() {
    this.http.get('assets/seed.sql', { responseType: 'text'})
    .subscribe(sql => {
      this.sqlitePorter.importSqlToDb(this.database, sql)
        .then(_ => {
          this.dbReady.next(true);
          console.log("WRITTEN SEED FILE");
        })
        .catch(e => console.error("there was an error", e));
    });
  }
 
  getDatabaseState() {
    return this.dbReady.asObservable();
  }

  addData(longitude, latitude, surface, tracking) {
    let data = [new Date(), longitude, latitude, surface, tracking];
    console.log("added data", data);
    return this.database.executeSql('INSERT INTO appdb (timestamp, longitude, latitude, surface, tracking) VALUES (?, ?, ?, ?, ?)', data);
  }

}
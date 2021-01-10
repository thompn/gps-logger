import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private url: string = "http://localhost:10000/api/v1/tracks"; //debugging locally - change to production server
  private headers: HttpHeaders = new HttpHeaders({
    "Content-Type": "application/json"
  });

  constructor(
    private http: HttpClient
  ) {

  }

  async updateDriverLocation(latitude: string, longitude: string, surface: string, tracking: string): Promise<any> {
    let body: any = {
      "timestamp": new Date(),
      "latitude": latitude.toString(),
      "longitude": longitude.toString(),
      "surface": surface,
      "tracking": tracking,
    };
    try {
      return this.http.post(`${this.url}`, body, { headers: this.headers, responseType: 'text'}).toPromise();
    }
    catch(err) {
      console.log(err.message)
    }
    
    }
}

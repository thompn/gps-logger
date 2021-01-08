import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private url: string = "https://api.sheetson.com/v2/sheets";
  private headers: HttpHeaders = new HttpHeaders({
    "X-Spreadsheet-Id": "1S6RM8eqqTn3BtZrEakTiytiNtnOoHs_cDON8-pjcwDc",
    "Authorization": "Bearer " + "-XNfXq7b4foSx7oAo-3Ufq9p9NQRN9HQI9MGRzPBoClqRi4ytecLoFEd0wU",
    "Content-Type": "application/json"
  });

  constructor(
    private http: HttpClient
  ) {

  }

  async updateDriverLocation(latitude: number, longitude: number, tracking: boolean, surface: string): Promise<any> {
    let body: any = {
      timestamp: new Date(),
      latitude: latitude,
      longitude: longitude,
      tracking: tracking,
      surface: surface,
    };
    return this.http.patch(`${this.url}/Location/2`, body, { headers: this.headers }).toPromise();
    console.log(body);
  }
}

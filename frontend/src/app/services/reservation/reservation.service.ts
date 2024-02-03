import { Injectable } from '@angular/core';
import { ConfigService } from '../config.service';
import { ApiService } from '../api.service';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ReservationService {

  constructor(
    private apiService: ApiService,
    private configService: ConfigService,
    private http: HttpClient,
    private router: Router
  ) {}


  getAvailability(id: string): Observable<any[]> {
    return this.http.get<any[]>(this.configService._getAvailability + id +"/availability");

  }
  postAvailability(availability: any): Observable<any[]> {
    return this.http.post<any[]>(this.configService._reservations_url + "/accomm/availability", availability);

  }

  getUserReservations(id: string): Observable<any[]> {
    return this.http.get<any[]>(this.configService._reservations_url +"/guest/" + id +"/reservations");

  }
  postReservation(availability: any): Observable<any[]> {
    return this.http.post<any[]>(this.configService._createReservation, availability);

  }
}

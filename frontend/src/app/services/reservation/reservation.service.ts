import { Injectable } from '@angular/core';
import { ConfigService } from '../config.service';
import { ApiService } from '../api.service';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ReservationService {
  constructor(
    private apiService: ApiService,
    private configService: ConfigService,
    private http: HttpClient,
    private router: Router
  ) {}

  getAvailability(id: string): Observable<any[]> {
    return this.http.get<any[]>(
      this.configService._getAvailability + id + '/availability'
    );
  }
  postAvailability(availability: any): Observable<any[]> {
    return this.http.post<any[]>(
      this.configService._reservations_url + '/accomm/availability',
      availability
    );
  }

  getUserReservations(id: string): Observable<any[]> {
    return this.http.get<any[]>(
      this.configService._reservations_url + '/guest/' + id + '/reservations'
    );
  }

  cancelReservations(id: string, headers: HttpHeaders): Observable<any[]> {
    const options = { headers };
    return this.http.delete<any[]>(
      this.configService._reservations_url + '/reservation/delete/' + id,
      options
    );
  }

  postReservation(availability: any): Observable<any[]> {
    return this.http.post<any[]>(
      this.configService._createReservation,
      availability
    );
  }
}

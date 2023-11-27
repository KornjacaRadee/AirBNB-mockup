import { Injectable } from '@angular/core';
import { ConfigService } from '../config.service';
import { ApiService } from '../api.service';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AccomodationService {
  constructor(
    private apiService: ApiService,
    private configService: ConfigService,
    private http: HttpClient,
    private router: Router
  ) {}

  getAccomodations(): Observable<any[]> {
    return this.http.get<any[]>(this.configService._accomodations_url + '/all');
  }

  searchAccomodations(
    searchTerm: string,
    minGuests: number,
    maxGuests: number
  ): Observable<any[]> {
    const body = {
      searchTerm: searchTerm,
      minGuests: minGuests.toString(),
      maxGuests: maxGuests.toString(),
    };

    return this.http.post<any[]>(
      this.configService._accomodations_url + '/search',
      body
    );
  }
}

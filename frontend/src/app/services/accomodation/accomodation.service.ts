import { Injectable } from '@angular/core';
import { ConfigService } from '../config.service';
import { ApiService } from '../api.service';
import { HttpClient, HttpHeaders } from '@angular/common/http';
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

  createAccommodation(
    headers: HttpHeaders,
    accommodation: any
  ): Observable<any[]> {
    const options = { headers };
    return this.http.post<any[]>(
      this.configService._accomodations_url + '/new',
      accommodation,
      options
    );
  }
  getUserAccommodations(headers: HttpHeaders): Observable<any[]> {
    const options = { headers };
    return this.http.get<any[]>(this.configService._userAccoms_url, options);
  }

  getAccommodation(id: string): Observable<any[]> {
    return this.http.get<any[]>(
      this.configService._accomodations_url + '/' + id
    );
  }

  deleteAccommodation(id: string,headers: HttpHeaders): Observable<any[]> {
    const options = { headers };
    return this.http.delete<any[]>(this.configService._accomodations_url + '/delete/' + id, options);
  }

  addAccommodationPictures(pictures: any): Observable<any[]> {

    return this.http.post<any[]>(this.configService._accomodations_url + '/accommodation/images', pictures);
  }

  getAccommodationPictures(id: any): Observable<any[]> {

    return this.http.get<any[]>(this.configService._accomodations_url +'/accommodation/' + id + '/images');
  }

  searchAccomodations(
    searchTerm: string,
    minGuests: number,
    maxGuests: number
  ): Observable<any[]> {
    const body = {
      location: searchTerm,
      minGuestNum: minGuests,
      maxGuestNum: maxGuests,
    };

    return this.http.post<any[]>(
      this.configService._accomodations_url + '/search',
      body
    );
  }
}

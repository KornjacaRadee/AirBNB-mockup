import { Injectable } from '@angular/core';
import { HttpClient,HttpHeaders } from '@angular/common/http';
import { ConfigService } from '../config.service';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';

interface Profile {
  email: string;
  username: string;
  firstName: string;
  lastName: string;
  address: string;
  role: string;
}

@Injectable({
  providedIn: 'root',
})
export class ProfilesService {
  constructor(
    private http: HttpClient,
    private configService: ConfigService,
    private router: Router
  ) {}

  getProfileByEmail(email: string): Observable<Profile> {
    return this.http.get<Profile>(
      `${this.configService._profiles_url}/u/${email}`
    );
  }

  getUserNotifications(headers: HttpHeaders): Observable<any> {
    const options = { headers };
    return this.http.get<any>(`${this.configService._notifications_url}/user-notifications`,options);
  }

  updateProfileByEmail(
    email: string,
    updatedProfile: Profile
  ): Observable<Profile> {
    return this.http.put<Profile>(
      `${this.configService._profiles_url}/update/${email}`,
      updatedProfile
    );
  }
}

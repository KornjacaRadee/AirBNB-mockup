import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ToastrModule } from 'ngx-toastr';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatToolbarModule } from '@angular/material/toolbar';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatFormFieldModule } from '@angular/material/form-field';
import { HomeComponent } from './home/home.component';
import { MatListModule } from '@angular/material/list';
import { MatCardModule } from '@angular/material/card';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { AuthService } from './services/auth.service';
import { ApiService } from './services/api.service';
import { RouterModule } from '@angular/router';
import { AccomodationService } from './services/accomodation/accomodation.service';
import { ReservationService } from './services/reservation/reservation.service';
import { DatePipe } from '@angular/common';

import {
  RECAPTCHA_SETTINGS,
  RecaptchaFormsModule,
  RecaptchaModule,
  RecaptchaSettings,
} from 'ng-recaptcha';
import { environment } from 'src/environments/environment';
import { PasswordRecoveryComponent } from './password-recovery/password-recovery.component';
import { NewPasswrodComponent } from './new-passwrod/new-passwrod.component';
import { NavbarComponent } from './navbar/navbar.component';
import { CreateAccommodationComponent } from './create-accommodation/create-accommodation.component';
import { ProfileComponent } from './profile/profile.component';
import { AccommodationPageComponent } from './accommodation-page/accommodation-page.component';
import { ProfilesService } from './services/profile/profiles.service';
import { CreatePeriodComponent } from './create-period/create-period.component';
import { UpdateProfileComponent } from './update-profile/update-profile.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    LoginComponent,
    RegisterComponent,
    PasswordRecoveryComponent,
    NewPasswrodComponent,
    NavbarComponent,
    CreateAccommodationComponent,
    ProfileComponent,
    AccommodationPageComponent,
    CreatePeriodComponent,
    UpdateProfileComponent,
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpClientModule,
    BrowserAnimationsModule,
    AppRoutingModule,
    MatToolbarModule,
    MatButtonModule,
    MatIconModule,
    MatInputModule,
    MatSidenavModule,
    MatFormFieldModule,
    MatListModule,
    MatCardModule,

    ReactiveFormsModule,
    FormsModule,
    RouterModule,
    RecaptchaModule,
    RecaptchaFormsModule,
  ],

  providers: [
    {
      provide: RECAPTCHA_SETTINGS,
      useValue: {
        siteKey: environment.recaptcha.siteKey,
      } as RecaptchaSettings,
    },
    DatePipe,
    AuthService,
    ApiService,
    AccomodationService,
    ReservationService,
    ProfilesService,
  ],

  bootstrap: [AppComponent],
})
export class AppModule {}

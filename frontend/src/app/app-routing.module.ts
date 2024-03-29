import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { RegisterComponent } from './register/register.component';
import { LoginComponent } from './login/login.component';
import { ProfileComponent } from './profile/profile.component';
import { PasswordRecoveryComponent } from './password-recovery/password-recovery.component';
import { NewPasswrodComponent } from './new-passwrod/new-passwrod.component';
import { NavbarComponent } from './navbar/navbar.component';
import { CreateAccommodationComponent } from './create-accommodation/create-accommodation.component';
import { AccommodationPageComponent } from './accommodation-page/accommodation-page.component';
import { CreatePeriodComponent } from './create-period/create-period.component';
import { BrowserModule } from '@angular/platform-browser';
import { CommonModule } from '@angular/common';
import { UpdateProfileComponent } from './update-profile/update-profile.component';
import { AddPicutresAccommComponent } from './add-picutres-accomm/add-picutres-accomm.component';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'home',
    pathMatch: 'full',
  },
  {
    path: 'home',
    component: HomeComponent,
  },
  {
    path: 'login',
    component: LoginComponent,
  },
  {
    path: 'register',
    component: RegisterComponent,
  },
  {
    path: 'profile',
    component: ProfileComponent,
  },
  {
    path: 'recovery',
    component: PasswordRecoveryComponent,
  },
  {
    path: 'new-password',
    component: NewPasswrodComponent,
  },
  {
    path: 'navbar',
    component: NavbarComponent,
  },
  {
    path: 'create-accommodation',
    component: CreateAccommodationComponent,
  },
  {
    path: 'accommodation-page',
    component: AccommodationPageComponent,
  },
  {
    path: 'create-period',
    component: CreatePeriodComponent,
  },
  {
    path: 'update-profile',
    component: UpdateProfileComponent,
  },
  {
    path: 'add-pictures-accomm',
    component: AddPicutresAccommComponent,
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes), BrowserModule, CommonModule],
  exports: [RouterModule],
})
export class AppRoutingModule {}

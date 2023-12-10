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
import { BrowserModule } from '@angular/platform-browser';
import { CommonModule } from '@angular/common';

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
];

@NgModule({
  imports: [RouterModule.forRoot(routes),BrowserModule,CommonModule],
  exports: [RouterModule],
})
export class AppRoutingModule {}

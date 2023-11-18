// register.component.ts
import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
})
export class RegisterComponent {

user: any;
registerUser() {
throw new Error('Method not implemented.');
}




  constructor(private authService: AuthService) {}


}

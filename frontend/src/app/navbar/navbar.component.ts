import { Component,Output, EventEmitter } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent {

  constructor(private authService: AuthService,private router: Router) { }
  signed(): boolean{
   return this.authService.isAuthenticated();
  }
  logout(){
    this.authService.logout();
  }

  isHost(): boolean{
    if(this.authService.getUserRole() == "host"){

      return true
    }else{
      return false
    }
  }
}

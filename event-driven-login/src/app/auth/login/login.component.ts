import { Component, inject, signal } from "@angular/core";
import { LoginService } from "./login.service";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { CommonModule } from "@angular/common";
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatIcon } from "@angular/material/icon";
import { MatCardModule } from "@angular/material/card";
import {MatButtonModule} from '@angular/material/button';
import { timer } from "rxjs";
import { MatProgressSpinner } from "@angular/material/progress-spinner";
import {
  trigger,
  transition,
  style,
  animate,
} from '@angular/animations';
import { RouterLink } from "@angular/router";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  imports: [CommonModule, ReactiveFormsModule, MatFormFieldModule,
    MatInputModule, MatIcon, MatCardModule,
    MatButtonModule, MatProgressSpinner, RouterLink],
  providers: [LoginService],
  animations: [
    trigger('fadeAnimation', [
      transition(':enter', [
        style({ opacity: 0 }),
        animate('200ms ease-in', style({ opacity: 1 }))
      ]),
      transition(':leave', [
        animate('200ms ease-out', style({ opacity: 0 }))
      ])
    ])
  ]
})
export class LoginComponent {
  private _loginService = inject(LoginService);
  private _formBuilder = inject(FormBuilder);
  form = this._formBuilder.group({
    username: ['', [Validators.required]],
    password: ['', [Validators.required]]
  });
  loginLoading = signal(false);
  loginSuccess = signal(false);
  loginError = signal(false);

  login() {
    const value = this.form.value;
    if (this.form.valid) {
      if (value.username && value.password) {
        this.loginLoading.set(true);
        this.loginSuccess.set(false);
        this.loginError.set(false);

        this._loginService.login(value.username, value.password).subscribe({
          next: () => {
            this.loginSuccess.set(true);
            console.log('Login successful');
          }, error: (error) => {
            this.loginError.set(true);
            console.error('Login failed', error);
            timer(1000).subscribe(() => {
              this.loginError.set(false);
            });
          }
        }).add(() => this.loginLoading.set(false));
      }
    }
  }
}

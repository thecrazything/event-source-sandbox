import { Routes } from '@angular/router';
import { LoginComponent } from './auth/login/login.component';
import { RegisterComponent } from './auth/register/register.component';
import { AuthComponent } from './auth/auth.component';

export const routes: Routes = [
  {
    path: '',
    redirectTo: '/auth/login',
    pathMatch: 'full'
  },
  {
    path: 'auth',
    component: AuthComponent,
    children: [
      {
        path: 'login',
        component: LoginComponent,
        data: { animation: 'LoginPage'}
      },
      {
        path: 'register',
        component: RegisterComponent,
        data: { animation: 'RegisterPage'}
      }
    ]
  },
  {
    path: 'login',
    redirectTo: '/auth/login',
    pathMatch: 'full',
    data: { animation: 'LoginPage'}
  },
  {
    path: 'register',
    redirectTo: '/auth/register',
    pathMatch: 'full',
    data: { animation: 'RegisterPage'}
  }
];

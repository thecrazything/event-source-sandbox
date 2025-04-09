import { inject, Injectable } from '@angular/core';
import { RequestService } from '../../service/request.service';
import { StatusMessage } from '../../service/status.message';
import { map, Observable } from 'rxjs';

@Injectable()
export class RegisterService {
  private _requestService = inject(RequestService);

  register(username: string, password: string): Observable<StatusMessage> {
    const payload = { username, password };
    return this._requestService.post<StatusMessage>('/api/register', payload).pipe(
      map(result => {
        if (result.status !== 'success') {
          throw new Error('Registration failed');
        }
        return result;
      })
    );
  }
}

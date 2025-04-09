import { inject, Injectable, signal } from "@angular/core";
import { map, mergeMap, Observable } from "rxjs";
import { RequestService } from "../../service/request.service";
import { StatusMessage } from "../../service/status.message";
import { HttpClient } from "@angular/common/http";

@Injectable()
export class LoginService {
  private _requestService = inject(RequestService);
  private _http = inject(HttpClient);

  login(username: string, password: string): Observable<StatusMessage> {
    const loginRequest = {
      username: username,
      password: password
    };
    return this._requestService.post<StatusMessage>('/api/login', loginRequest).pipe(mergeMap(result => {
      if (result.status === 'success') {
        // if the logn was succesful, we need to fetch the session cookie
        return this._http.get<void>('/api/session').pipe(map(() => result));
      } else {
        return new Observable<StatusMessage>(observer => {
          observer.error(result);
          observer.complete();
        });
      }
    }))
  }
}

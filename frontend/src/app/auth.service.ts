import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { BehaviorSubject, combineLatest, of } from 'rxjs';
import { catchError, map, switchMap, tap } from 'rxjs/operators';
import { environment } from '../environments/environment';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private authed$ = new BehaviorSubject<boolean>(false);
  private roles$ = new BehaviorSubject<string[]>([]);
  private user$ = new BehaviorSubject<{ id: number; email: string; name: string } | null>(null);

  constructor(private http: HttpClient) {}

  private headers(): HttpHeaders {
    let h = new HttpHeaders({ 'Content-Type': 'application/json' });
    const token = localStorage.getItem('accessToken');
    if (token) { h = h.set('Authorization', `Bearer ${token}`); }
    return h;
  }

  isAuthed() { return this.authed$.asObservable(); }
  roles() { return this.roles$.asObservable(); }
  user() { return this.user$.asObservable(); }

  check() {
    // Verify, then try refresh if invalid, then verify again
    const verify$ = () => this.http.get<{ success: boolean; data: any }>(`${environment.apiBase}/v1/auth/verify`, { headers: this.headers(), withCredentials: true })
      .pipe(catchError(() => of({ success: true, data: { valid: false } } as any)));

    verify$().pipe(
      switchMap(first => {
        if (first?.data?.valid) return of(first);
        // attempt refresh then verify again
        return this.http.post(`${environment.apiBase}/v1/auth/refresh`, {}, { headers: this.headers(), withCredentials: true })
          .pipe(
            catchError(() => of(null)),
            switchMap(() => verify$())
          );
      }),
      tap(res => this.authed$.next(!!res?.data?.valid)),
      switchMap(res => {
        if (!res?.data?.valid) {
          this.user$.next(null);
          this.roles$.next([]);
          return of(null);
        }
        return combineLatest([
          this.http.get<{ success: boolean; data: any }>(`${environment.apiBase}/v1/me`, { headers: this.headers(), withCredentials: true }).pipe(catchError(() => of({ success: false } as any))),
          this.http.get<{ success: boolean; data: string[] }>(`${environment.apiBase}/v1/auth/roles`, { headers: this.headers(), withCredentials: true }).pipe(catchError(() => of({ success: false, data: [] } as any))),
        ]).pipe(map(([me, roles]) => ({ me, roles })));
      })
    ).subscribe((combo: any) => {
      if (!combo) return;
      const me = combo.me?.data;
      const roles = combo.roles?.data || [];
      this.user$.next(me || null);
      this.roles$.next(Array.isArray(roles) ? roles : []);
    });
  }
}

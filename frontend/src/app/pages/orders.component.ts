import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ApiService } from '../api.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-orders',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './orders.component.html'
})
export class OrdersPageComponent {
  loading = false;
  items: any[] = [];
  constructor(private api: ApiService, private router: Router) { this.refresh(); }
  refresh() {
    this.loading = true;
    this.api.listMyPurchases().subscribe({
      next: res => { this.items = res.data || []; this.loading = false; },
      error: () => { this.items = []; this.loading = false; }
    });
  }
  read(bookId: string) { this.router.navigate(['/book', bookId]); }
}

import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-authors',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './authors.component.html'
})
export class AuthorsPageComponent implements OnInit {
  authors: any[] = [];
  loading = false;
  q = '';
  genre = '';
  publishers: any[] = [];
  selectedPublisherId: string | null = null;
  inviteMessage = '';

  constructor(private api: ApiService) {}

  ngOnInit(): void {
    this.loadPublishers();
    this.search();
  }

  loadPublishers() {
    this.api.listMyPublishers().subscribe({
      next: res => {
        this.publishers = res.data || [];
        if (!this.selectedPublisherId && this.publishers.length) {
          this.selectedPublisherId = this.publishers[0].id;
        }
      },
      error: () => { this.publishers = []; }
    });
  }

  search() {
    this.loading = true;
    this.api.getPublicAuthors({ q: this.q, genre: this.genre }).subscribe({
      next: res => { this.authors = res.data || []; this.loading = false; },
      error: () => { this.authors = []; this.loading = false; }
    });
  }

  invite(author: any) {
    if (!this.selectedPublisherId) {
      this.inviteMessage = 'Select a publisher profile first.';
      return;
    }
    this.inviteMessage = 'Sending invite...';
    this.api.createCollaboration(this.selectedPublisherId, author.id).subscribe({
      next: () => { this.inviteMessage = 'Collaboration request sent.'; },
      error: () => { this.inviteMessage = 'Unable to send request.'; }
    });
  }
}

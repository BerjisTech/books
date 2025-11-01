import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BooksListComponent } from '../books-list/books-list.component';
import { BooksService } from 'src/app/services/books.service';
import { ApiService } from '../../api.service';
import { Book } from '../../interfaces/book';

@Component({
  selector: 'app-books',
  standalone: true,
  templateUrl: './books.component.html',
  styleUrl: './books.component.scss',
  imports: [CommonModule, BooksListComponent]
})
export class BooksComponent implements OnInit {
  books: Book[] = [];
  loading = false;

  constructor(private api: ApiService, private booksService: BooksService) {}

  ngOnInit(): void {
    this.load();
  }

  load() {
    this.loading = true;
    this.api.getPublicBooks().subscribe({
      next: res => {
        const data = res?.data || [];
        const fallbacks = this.booksService.getDummyBooks(data.length > 0 ? data.length : 5);
        this.books = data.map((b: any, idx: number) => ({
          title: b.title,
          slag: b.slug || b.id,
          cover: b.coverImageUrl || fallbacks[idx % fallbacks.length].cover,
          short_description: b.description,
          author_name: b.authorName,
          price: b.priceAmount,
          priceCurrency: b.priceCurrency
        }));
        this.loading = false;
      },
      error: () => {
        this.books = this.booksService.getDummyBooks(8);
        this.loading = false;
      }
    });
  }
}

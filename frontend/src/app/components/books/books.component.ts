import { Component } from '@angular/core';
import { Book } from '../../interfaces/book';
import { BooksService } from 'src/app/services/books.service';
import { BooksListComponent } from "../books-list/books-list.component";
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-books',
  standalone: true,
  templateUrl: './books.component.html',
  styleUrl: './books.component.scss',
  imports: [CommonModule, BooksListComponent]
})
export class BooksComponent {
  constructor(private booksService: BooksService) { }
  dummyBooks(bookCount: number = 8): Book[] {
    return this.booksService.getDummyBooks(bookCount);
  }
}
